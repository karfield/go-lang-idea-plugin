package ro.redeul.google.go.completion;

import java.io.IOException;

public class GoCompletionBugsTestCase extends GoCompletionTestCase{
    protected String getTestDataRelativePath() {
        return super.getTestDataRelativePath() + "bugs";
    }

    public void testGH218_MissingTypeSpec() throws IOException {
        _testVariants();
    }

    public void testGH530_MissingFunctionName() throws IOException {
        _testVariants();
    }

}
